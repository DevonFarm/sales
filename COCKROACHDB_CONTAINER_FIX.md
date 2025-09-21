# ✅ CockroachDB Container Initialization Fix

## Problem Fixed

The integration and e2e test jobs were failing with this error:
```log
/cockroach/cockroach.sh: line 278: 1: error: mode unset, can be shell, bash, or cockroach command   (start-single-node, sql, etc.)
```

## Root Cause

The error occurred because:
1. **GitHub Actions Services Limitation**: The `services` configuration in GitHub Actions doesn't easily allow overriding the container's default command/entrypoint
2. **CockroachDB Requirement**: The CockroachDB container requires a specific command (`start-single-node`) to be passed when starting
3. **Missing Command**: The service was starting without the required `start-single-node` command, causing the "mode unset" error

## Solution Applied

**Replaced GitHub Actions services with manual Docker container steps:**

### Before (Problematic):
```yaml
services:
  cockroachdb:
    image: cockroachdb/cockroach:latest
    ports:
      - 26257:26257
    options: --health-cmd="curl -f http://localhost:8080/health" --health-interval=10s --health-timeout=5s --health-retries=5
```

### After (Fixed):
```yaml
steps:
- name: Start CockroachDB
  run: |
    docker run -d --name cockroach \
      -p 26257:26257 \
      -p 8080:8080 \
      cockroachdb/cockroach:latest \
      start-single-node --insecure --listen-addr=0.0.0.0:26257 --http-addr=0.0.0.0:8080

- name: Wait for CockroachDB
  run: |
    timeout 60 bash -c 'until curl -f http://localhost:8080/health; do sleep 2; done'
```

## ✅ Jobs Updated

**Both jobs that use CockroachDB were fixed:**

### 1. **integration-tests** Job
- ✅ Replaced `services` with manual Docker container step
- ✅ Added proper `start-single-node` command
- ✅ Added health check wait step
- ✅ Moved wait step before Go setup for proper sequencing

### 2. **e2e-tests** Job  
- ✅ Replaced `services` with manual Docker container step
- ✅ Added proper `start-single-node` command
- ✅ Added health check wait step
- ✅ Moved wait step before Go/Node setup for proper sequencing

## 🎯 Key Improvements

1. **Explicit Command**: CockroachDB now starts with the required `start-single-node --insecure` command
2. **Proper Networking**: Added both ports 26257 (SQL) and 8080 (HTTP admin) 
3. **Health Checking**: Maintained the health check but moved it to a separate step
4. **Better Sequencing**: Database starts and becomes ready before other setup steps
5. **More Control**: Full control over container startup parameters

## 🚀 Expected Results

- ✅ CockroachDB containers will start successfully
- ✅ No more "mode unset" errors
- ✅ Database will be ready for migrations and tests
- ✅ Integration tests can connect to the database
- ✅ E2E tests can connect to the database

## 📝 Technical Details

**Container Command:**
```bash
docker run -d --name cockroach \
  -p 26257:26257 \        # SQL port
  -p 8080:8080 \          # HTTP admin port
  cockroachdb/cockroach:latest \
  start-single-node --insecure --listen-addr=0.0.0.0:26257 --http-addr=0.0.0.0:8080
```

**Health Check:**
```bash
timeout 60 bash -c 'until curl -f http://localhost:8080/health; do sleep 2; done'
```

This approach provides more reliability and control compared to GitHub Actions services for containers that need specific startup commands.

The CockroachDB initialization error should now be completely resolved! 🎉