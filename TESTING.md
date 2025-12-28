# Testing Documentation

This document describes the testing infrastructure for taskspec-go.

## Test Types

### 1. Unit Tests

Standard Go unit tests that verify the parser functionality with known inputs and expected outputs.

**Run:**
```bash
go test -v
```

**Coverage:**
```bash
go test -cover -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### 2. Fuzz Tests

Fuzz tests use Go's built-in fuzzing capabilities to test the parser with random inputs, helping to discover edge cases, panics, and crashes.

**Files:**
- `parser_fuzz_test.go` - Contains fuzz tests for `Parse()` and `ParseLines()` functions

**Run:**
```bash
# Run for a specific duration
go test -fuzz=FuzzParser -fuzztime=2m
go test -fuzz=FuzzParseLines -fuzztime=2m

# Run all fuzz tests quickly (seed corpus only)
go test -v
```

**What it tests:**
- The parser should never panic or crash, regardless of input
- Random combinations of valid and invalid taskspec syntax
- Edge cases with empty strings, special characters, and malformed metadata

### 3. ABNF Grammar Tests

Tests the parser against samples generated from the taskspec ABNF grammar specification.

**Files:**
- `parser_abnfgen_test.go` - Test runner for grammar-generated samples
- `testdata/grammar.abnf` - The official taskspec ABNF grammar
- `testdata/generate_samples.py` - Python script to generate random taskspec entries

**Generate samples:**
```bash
# Generate samples for 60 seconds
python3 testdata/generate_samples.py 60 > testdata/abnfgen_samples.txt

# Generate samples for 10 minutes (default)
python3 testdata/generate_samples.py 600 > testdata/abnfgen_samples.txt
```

**Run tests:**
```bash
# First generate samples, then run the test
python3 testdata/generate_samples.py 10 > testdata/abnfgen_samples.txt
go test -v -run TestAbnfGenSamples
```

**What it tests:**
- Parser robustness against a wide variety of valid taskspec entries
- Ensures the parser doesn't panic on grammar-compliant inputs
- Validates that the parser can handle the full range of taskspec features

## Continuous Integration

### GitHub Actions Workflow

The repository includes a comprehensive GitHub Actions workflow (`.github/workflows/advanced-testing.yml`) that runs:

1. **Fuzz Testing** - Runs fuzz tests for 4 minutes total (2 minutes per fuzz function)
2. **ABNF Grammar Testing** - Generates samples for 10 minutes and tests them
3. **Regular Unit Tests** - Runs all standard tests with race detection

**Triggers:**
- On every push to the `main` branch
- On a daily schedule at 00:00 UTC
- Manual workflow dispatch

**Failure Handling:**
- Automatically creates a GitHub issue when tests fail
- Includes run details, commit information, and links to the workflow run
- Uploads test artifacts (coverage, generated samples) for investigation

## Running All Tests

To run the complete test suite locally:

```bash
# 1. Run regular tests
go test -v -race -coverprofile=coverage.out

# 2. Run fuzz tests (adjust time as needed)
go test -fuzz=FuzzParser -fuzztime=30s
go test -fuzz=FuzzParseLines -fuzztime=30s

# 3. Generate and run ABNF grammar tests
python3 testdata/generate_samples.py 60 > testdata/abnfgen_samples.txt
go test -v -run TestAbnfGenSamples

# 4. Clean up generated files
rm testdata/abnfgen_samples.txt
```

## Test Statistics

The ABNF grammar test reports statistics about:
- Total samples tested
- Successfully parsed samples
- Parse failures (allowed - grammar may generate semantically invalid but syntactically valid input)
- Panics (not allowed - test fails if any occur)

## Contributing

When adding new features to the parser:

1. Add unit tests for the new functionality
2. Update fuzz test seed corpus if applicable
3. Ensure all tests pass before submitting a PR
4. The CI workflow will automatically run all tests

## Troubleshooting

**Issue: Fuzz test finds a crash**
- The fuzzer saves failing inputs in `testdata/fuzz/`
- Run `go test -fuzz=FuzzParser` to reproduce
- Fix the crash and verify with the saved input

**Issue: ABNF test takes too long**
- Generate fewer samples by reducing the time parameter
- Consider testing with a smaller sample set locally

**Issue: CI creates duplicate issues**
- Issues are created only on push to main or scheduled runs
- Check if the issue already exists before re-running
