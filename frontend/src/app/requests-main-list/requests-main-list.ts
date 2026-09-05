import { Component, inject } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { WorkspaceStateService } from '../core/services/workspace-state.service';

@Component({
  selector: 'app-requests-main-list',
  imports: [FormsModule],
  templateUrl: './requests-main-list.html',
  styleUrl: './requests-main-list.scss',
  host: {
    class: 'main-column',
    'aria-label': 'Requests',
  },
})
export class RequestsMainList {
  protected readonly state = inject(WorkspaceStateService);
}
