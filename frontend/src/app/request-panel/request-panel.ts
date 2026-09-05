import { Component, inject } from '@angular/core';
import { WorkspaceStateService } from '../core/services/workspace-state.service';
import { RequestEditor } from '../request-editor/request-editor';
import { ResponseViewer } from '../response-viewer/response-viewer';

@Component({
  selector: 'app-request-panel',
  imports: [RequestEditor, ResponseViewer],
  templateUrl: './request-panel.html',
  styleUrl: './request-panel.scss',
  host: {
    class: 'sidebar sidebar-right',
    'aria-label': 'Request and response panel',
  },
})
export class RequestPanel {
  protected readonly state = inject(WorkspaceStateService);
}
