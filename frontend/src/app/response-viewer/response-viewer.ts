import { Component, inject } from '@angular/core';
import { WorkspaceStateService } from '../core/services/workspace-state.service';

@Component({
  selector: 'app-response-viewer',
  templateUrl: './response-viewer.html',
  styleUrl: './response-viewer.scss',
})
export class ResponseViewer {
  protected readonly state = inject(WorkspaceStateService);
}
