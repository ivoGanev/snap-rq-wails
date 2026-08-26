import { AfterViewInit, Component, effect, inject, OnInit, signal } from '@angular/core';
import { WML } from '@wailsio/runtime';
import { WailsService } from './wails.service';
import { RequestApiService, type HttpRequest, type HttpResponse } from './services/request.service';
import { CollectionApiService, type Collection } from './services/collection.service';

@Component({
  selector: 'app-root',
  imports: [],
  templateUrl: './app.html',
  styleUrl: './app.scss',
})
export class App implements OnInit, AfterViewInit {
  private readonly wails = inject(WailsService);
  private readonly requestApi = inject(RequestApiService);
  private readonly collectionApi = inject(CollectionApiService);

  readonly currentTime = this.wails.currentTime;
  readonly collections = this.collectionApi.collections;
  readonly requests = this.requestApi.requests;
  readonly responses = this.requestApi.responses;
  readonly loading = signal(false);
  readonly error = signal<string | null>(null);
  readonly selectedCollection = signal<Collection | null>(null);
  readonly selectedRequest = signal<HttpRequest | null>(null);
  readonly selectedResponse = signal<HttpResponse | null>(null);

  constructor() {
    effect(() => {
      const req = this.selectedRequest();
      if (req) {
        this.loadResponses(req.id);
      } else {
        this.requestApi.responses.set([]);
        this.selectedResponse.set(null);
      }
    });
  }

  async ngOnInit(): Promise<void> {
    await this.loadCollections();
  }

  ngAfterViewInit(): void {
    WML.Enable();
  }

  async loadCollections(): Promise<void> {
    this.loading.set(true);
    this.error.set(null);
    try {
      await this.collectionApi.loadAll();
      const all = this.collections();
      if (all.length > 0) {
        this.selectCollection(all[0]);
      }
    } catch (err) {
      console.error(err);
      this.error.set('Failed to load collections.');
    } finally {
      this.loading.set(false);
    }
  }

  async selectCollection(collection: Collection): Promise<void> {
    this.selectedCollection.set(collection);
    this.selectedRequest.set(null);
    this.selectedResponse.set(null);
    try {
      await this.requestApi.loadForCollection(collection.id);
    } catch (err) {
      console.error(err);
    }
  }

  async loadResponses(requestId: number): Promise<void> {
    try {
      await this.requestApi.loadResponsesForRequest(requestId);
      const all = this.responses();
      this.selectedResponse.set(all.length > 0 ? all[0] : null);
    } catch (err) {
      console.error(err);
    }
  }

  selectRequest(req: HttpRequest): void {
    this.selectedRequest.set(req);
  }

  selectResponse(resp: HttpResponse): void {
    this.selectedResponse.set(resp);
  }
}
