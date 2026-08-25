import { AfterViewInit, Component, inject, OnInit, signal } from '@angular/core';
import { WML } from '@wailsio/runtime';
import { WailsService } from './wails.service';
import { RequestApiService, type HttpRequest } from './core/services/request.service';

@Component({
  selector: 'app-root',
  imports: [],
  templateUrl: './app.html',
  styleUrl: './app.scss',
})
export class App implements OnInit, AfterViewInit {
  private readonly wails = inject(WailsService);
  private readonly requestApi = inject(RequestApiService);

  readonly currentTime = this.wails.currentTime;
  readonly requests = this.requestApi.requests;
  readonly loading = signal(false);
  readonly error = signal<string | null>(null);
  readonly selectedRequest = signal<HttpRequest | null>(null);

  async ngOnInit(): Promise<void> {
    await this.loadRequests();
  }

  ngAfterViewInit(): void {
    WML.Enable();
  }

  async loadRequests(): Promise<void> {
    this.loading.set(true);
    this.error.set(null);
    try {
      await this.requestApi.loadAll();
    } catch (err) {
      console.error(err);
      this.error.set('Failed to load requests.');
    } finally {
      this.loading.set(false);
    }
  }

  selectRequest(req: HttpRequest): void {
    this.selectedRequest.set(req);
  }
}
