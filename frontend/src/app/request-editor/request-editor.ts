import { Component, inject } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { WorkspaceStateService } from '../core/services/workspace-state.service';

@Component({
  selector: 'app-request-editor',
  imports: [FormsModule],
  templateUrl: './request-editor.html',
  styleUrl: './request-editor.scss',
})
export class RequestEditor {
  protected readonly state = inject(WorkspaceStateService);
}
