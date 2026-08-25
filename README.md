# Snap RQ

Your API testing buddy!

## How to run in dev mode (Wails)

1. Navigate to your project directory in the terminal.

2. To run your application in development mode, use the following command:

   ```
   wails3 dev
   ```

   This will start your application and enable hot-reloading for both frontend and backend changes.

3. To build your application for production, use:

   ```
   wails3 build
   ```

   This will create a production-ready executable in the `build` directory.

### Rebuilding bindings
- wails3 generate bindings -ts -i -clean=true

## Seeding mock data

### 50 mixed-size requests in the default app database
```
python scripts/mock_requests.py
```

### 200 requests in a custom file
```
python scripts/mock_requests.py --count 200 --db ./test.db
```

### Small data only
```
python scripts/mock_requests.py --count 10 --min-len 5 --max-len 30
```

### Large data for UI stress testing
```
python scripts/mock_requests.py --count 25 --min-len 500 --max-len 4000
```

### Reproducible data + clear existing rows first
```
python scripts/mock_requests.py --count 100 --seed 42 --clear
```