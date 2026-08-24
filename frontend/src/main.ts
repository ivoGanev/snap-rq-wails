import { bootstrapApplication } from '@angular/platform-browser';
import { appConfig } from './app/app.config';
import { App } from './app/app';

// Import the generated event bindings so that the typed `time` event
// declared in Go is registered with the runtime.
import '../bindings/github.com/wailsapp/wails/v3/internal/eventcreate';

bootstrapApplication(App, appConfig).catch((err) => console.error(err));
