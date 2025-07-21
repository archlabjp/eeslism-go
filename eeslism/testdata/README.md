# Test Data for EESLISM PCM Module

This directory contains test data files for PCM (Phase Change Material) CHARTABLE functionality.

## Files

### pcm_enthalpy_test.txt
Test data for PCM enthalpy table testing.
- Format: `temperature enthalpy`
- Temperature range: 20°C to 32°C
- Contains 7 data points with a phase change peak at 26-28°C

### pcm_conductivity_test.txt
Test data for PCM thermal conductivity table testing.
- Format: `temperature conductivity`
- Temperature range: 20°C to 32°C
- Linear increase from 0.15 to 0.21 W/mK

## Usage

These files are used by the following test functions in `blpcm_test.go`:
- `TestPCMdata` (CHARTABLE reading subtest)
- `TestTableRead`
- `TestFNPCMenthalpy_table_lib`

## Data Format

Both files use the same format:
```
temperature_value characteristic_value
```

Where:
- `temperature_value`: Temperature in Celsius
- `characteristic_value`: Either enthalpy (J/m³) or thermal conductivity (W/mK)