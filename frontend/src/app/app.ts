import { AfterViewInit, Component, inject, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { WML } from '@wailsio/runtime';
import { WailsService } from './wails.service';

@Component({
  selector: 'app-root',
  imports: [FormsModule],
  templateUrl: './app.html',
  styleUrl: './app.scss',
})
export class App implements AfterViewInit {
  private readonly wails = inject(WailsService);

  readonly name = signal('');
  readonly greeting = signal<string | null>(null);
  readonly currentTime = this.wails.currentTime;

  ngAfterViewInit(): void {
    // Enable Wails Markup Language handlers for data-wml-openURL links.
    WML.Enable();
  }

  async greet(): Promise<void> {
    try {
      const message = await this.wails.greet(this.name());
      this.greeting.set(message);
    } catch (err) {
      console.error(err);
    }
  }

  dismiss(): void {
    this.greeting.set(null);
  }
}
