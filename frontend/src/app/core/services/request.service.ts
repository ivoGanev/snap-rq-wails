import { Injectable, signal } from '@angular/core';
import * as RequestService from '../../../../bindings/snap-rq/services';
import type { HttpRequest } from '../../../../bindings/snap-rq/models';

export type { HttpRequest };

/**
 * Angular wrapper around the Wails-generated RequestService bindings.
 */
@Injectable({ providedIn: 'root' })
export class RequestApiService {
  readonly requests = signal<HttpRequest[]>([]);

  async create(req: Omit<HttpRequest, 'id'>): Promise<HttpRequest> {
    const created = await RequestService.RequestService.CreateRequest(req as HttpRequest);
    await this.loadAll();
    return created;
  }

  async get(id: number): Promise<HttpRequest> {
    return RequestService.RequestService.GetRequest(id);
  }

  async loadAll(): Promise<void> {
    const all = await RequestService.RequestService.GetAllRequests();
    this.requests.set(all ?? []);
  }

  async update(req: HttpRequest): Promise<HttpRequest> {
    const updated = await RequestService.RequestService.UpdateRequest(req);
    await this.loadAll();
    return updated;
  }

  async delete(id: number): Promise<void> {
    await RequestService.RequestService.DeleteRequest(id);
    await this.loadAll();
  }
}
