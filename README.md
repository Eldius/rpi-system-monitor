# rpi-system-monitor

## references

### Voltage measurements

```shell
# core voltage
vcgencmd measure_volts core
```

**Options for vcgencmd measure_volts:**
  - core: Measures the voltage of the VideoCore (CPU).
  - sdram_c: Measures the voltage of the RAM controller.
  - sdram_i: Measures the RAM I/O voltage.
  - sdram_p: Measures the RAM Phy voltage.

