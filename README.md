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

### Scale test script
```
python scripts/mock_requests.py
```

### Reproducible data + clear existing rows first
```
python scripts/mock_requests.py --count 100 --seed 42 --clear
```