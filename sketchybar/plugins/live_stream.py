#!/usr/bin/env python3
"""
Live Streaming Daemon for SketchyBar (Watching Mode)
Streams real-time BTC and Gold ticks via WebSocket and samples system stats via Mach Kernel.
"""

import asyncio
import ctypes
import ctypes.util
import json
import os
import subprocess
import sys
import time

# --- Setup Mach Kernel for Zero-Overhead System Stats ---
libc = ctypes.CDLL(ctypes.util.find_library("c"))

class HostCpuLoadInfo(ctypes.Structure):
    _fields_ = [("cpu_ticks", ctypes.c_uint * 4)]

HOST_CPU_LOAD_INFO = 3
HOST_CPU_LOAD_INFO_COUNT = ctypes.c_uint(ctypes.sizeof(HostCpuLoadInfo) // ctypes.sizeof(ctypes.c_int))

def get_cpu_ticks():
    try:
        host = libc.mach_host_self()
        info = HostCpuLoadInfo()
        count = HOST_CPU_LOAD_INFO_COUNT
        libc.host_statistics(host, HOST_CPU_LOAD_INFO, ctypes.byref(info), ctypes.byref(count))
        return list(info.cpu_ticks)
    except Exception:
        return [0, 0, 0, 0]

def get_ram_usage():
    try:
        # Fast query via top
        out = subprocess.check_output(
            ["top", "-l", "1", "-n", "0"],
            stderr=subprocess.DEVNULL,
            timeout=1
        ).decode()
        for line in out.splitlines():
            if "PhysMem:" in line:
                parts = line.split()
                # e.g., PhysMem: 12G used (1750M wired, 2127M compressor), 3865M unused.
                u_str = parts[1].replace("G", "").replace("M", "")
                u_val = float(u_str)
                if "M" in parts[1]:
                    u_val /= 1024.0
                
                # find unused
                for i, p in enumerate(parts):
                    if "unused" in p:
                        f_str = parts[i-1].replace("G", "").replace("M", "")
                        f_val = float(f_str)
                        if "M" in parts[i-1]:
                            f_val /= 1024.0
                        tot = u_val + f_val
                        if tot > 0:
                            return f"{int(round((u_val / tot) * 100))}%"
    except Exception:
        pass
    return "0%"

def get_disk_usage():
    try:
        out = subprocess.check_output(["df", "-H", "/"], stderr=subprocess.DEVNULL).decode()
        lines = out.strip().splitlines()
        if len(lines) >= 2:
            return lines[1].split()[4]
    except Exception:
        pass
    return "0%"

# --- State ---
last_btc_label = ""
last_btc_color = ""
last_gold_label = ""
last_gold_color = ""

def format_candle(candle):
    open_p = float(candle["o"])
    curr_p = float(candle["c"])
    diff = curr_p - open_p
    pct = (diff / open_p) * 100
    arrow = "▲" if diff >= 0 else "▼"
    color = "0xffa6e3a1" if diff >= 0 else "0xfff38ba8"
    sign = "+" if pct >= 0 else ""
    formatted = f"${int(round(curr_p)):,} {sign}{pct:.2f}% {arrow}"
    return formatted, color

async def stream_crypto():
    global last_btc_label, last_btc_color, last_gold_label, last_gold_color
    import websockets
    uri = "wss://stream.binance.com:9443/stream?streams=btcusdt@kline_15m/paxgusdt@kline_15m"

    while True:
        try:
            async with websockets.connect(uri, ping_interval=20, ping_timeout=10) as ws:
                async for message in ws:
                    data = json.loads(message)
                    stream = data.get("stream", "")
                    kline = data.get("data", {}).get("k", {})
                    if not kline:
                        continue

                    if "btcusdt" in stream:
                        label, color = format_candle(kline)
                        if label != last_btc_label or color != last_btc_color:
                            last_btc_label = label
                            last_btc_color = color
                            subprocess.run(["sketchybar", "--set", "btc", f"label={label}", f"label.color={color}"])

                    elif "paxgusdt" in stream:
                        label, color = format_candle(kline)
                        if label != last_gold_label or color != last_gold_color:
                            last_gold_label = label
                            last_gold_color = color
                            subprocess.run(["sketchybar", "--set", "gold", f"label={label}", f"label.color={color}"])
        except Exception:
            await asyncio.sleep(2)

async def stream_stats():
    last_ticks = get_cpu_ticks()
    disk_counter = 0
    disk_val = get_disk_usage()
    ram_val = get_ram_usage()

    while True:
        await asyncio.sleep(1.0)
        curr_ticks = get_cpu_ticks()
        deltas = [b - a for a, b in zip(last_ticks, curr_ticks)]
        total = sum(deltas)
        last_ticks = curr_ticks

        if total > 0:
            user, sys_t, idle, nice = deltas
            used = (user + sys_t + nice) / total * 100
            cpu_val = f"{int(round(used))}%"
        else:
            cpu_val = "0%"

        disk_counter += 1
        if disk_counter >= 15:
            disk_counter = 0
            disk_val = get_disk_usage()
            ram_val = get_ram_usage()

        subprocess.run([
            "sketchybar",
            "--set", "cpu", f"label={cpu_val}",
            "--set", "ram", f"label={ram_val}",
            "--set", "disk", f"label={disk_val}"
        ])

async def main():
    await asyncio.gather(stream_crypto(), stream_stats())

if __name__ == "__main__":
    try:
        asyncio.run(main())
    except KeyboardInterrupt:
        pass
