import { AfterViewInit, Component, effect, inject, OnInit, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { WML } from '@wailsio/runtime';
import { WailsService } from './wails.service';
import { RequestApiService, type HttpRequest, type HttpResponse } from './core/services/request.service';
import { CollectionApiService, type Collection } from './core/services/collection.service';

type RightPanelMode = 'response' | 'edit';

@Component({
  selector: 'app-root',
  imports: [FormsModule],
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
  readonly rightPanelMode = signal<RightPanelMode>('response');
  readonly draftRequest = signal<HttpRequest | null>(null);

  constructor() {
    effect(() => {
      const req = this.selectedRequest();
      if (req) {
        this.loadResponses(req.id);
        this.draftRequest.set({ ...req });
      } else {
        this.requestApi.responses.set([]);
        this.selectedResponse.set(null);
        this.draftRequest.set(null);
        this.rightPanelMode.set('response');
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
    this.rightPanelMode.set('response');
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
    this.rightPanelMode.set('response');
  }

  selectResponse(resp: HttpResponse): void {
    this.selectedResponse.set(resp);
  }

  setRightPanelMode(mode: RightPanelMode): void {
    this.rightPanelMode.set(mode);
    if (mode === 'edit') {
      const req = this.selectedRequest();
      this.draftRequest.set(req ? { ...req } : null);
    }
  }

  updateDraft<K extends keyof HttpRequest>(field: K, value: HttpRequest[K]): void {
    const draft = this.draftRequest();
    if (!draft) return;
    this.draftRequest.set({ ...draft, [field]: value });
  }

  async saveRequest(): Promise<void> {
    const draft = this.draftRequest();
    const collection = this.selectedCollection();
    if (!draft || !collection) return;

    this.loading.set(true);
    this.error.set(null);
    try {
      const updated = await this.requestApi.update({ ...draft, collection_id: collection.id });
      this.selectedRequest.set(updated);
      await this.requestApi.loadForCollection(collection.id);
      this.rightPanelMode.set('response');
    } catch (err) {
      console.error(err);
      this.error.set('Failed to save request.');
    } finally {
      this.loading.set(false);
    }
  }

  cancelEdit(): void {
    const req = this.selectedRequest();
    this.draftRequest.set(req ? { ...req } : null);
    this.rightPanelMode.set('response');
  }

  async sendRequest(req: HttpRequest, event: MouseEvent): Promise<void> {
    event.stopPropagation();
    this.loading.set(true);
    try {
      const resp = await this.requestApi.execute(req.id);
      if (this.selectedCollection()) {
        await this.requestApi.loadForCollection(this.selectedCollection()!.id);
      }
      if (this.selectedRequest()?.id === req.id) {
        await this.loadResponses(req.id);
        this.selectedResponse.set(resp);
        this.rightPanelMode.set('response');
      }
    } catch (err) {
      console.error(err);
      this.error.set('Failed to send request.');
    } finally {
      this.loading.set(false);
    }
  }
}
