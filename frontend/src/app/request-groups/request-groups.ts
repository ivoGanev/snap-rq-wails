import { Component, inject } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { WorkspaceStateService } from '../core/services/workspace-state.service';

@Component({
  selector: 'app-request-groups',
  imports: [FormsModule],
  templateUrl: './request-groups.html',
  styleUrl: './request-groups.scss',
  host: {
    class: 'sidebar sidebar-left',
    'aria-label': 'Collections and Favourites',
  },
})
export class RequestGroups {
  protected readonly state = inject(WorkspaceStateService);
}
