import { AfterViewInit, Component, inject, OnInit, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { WML } from '@wailsio/runtime';
import { WorkspaceStateService } from './core/services/workspace-state.service';
import { RequestGroups } from './request-groups/request-groups';
import { RequestsMainList } from './requests-main-list/requests-main-list';
import { RequestPanel } from './request-panel/request-panel';

@Component({
  selector: 'app-root',
  imports: [FormsModule, RequestGroups, RequestsMainList, RequestPanel],
  templateUrl: './app.html',
  styleUrl: './app.scss',
})
export class App implements OnInit, AfterViewInit {
  protected readonly state = inject(WorkspaceStateService);

  readonly leftColumnWidth = signal(220);
  readonly rightColumnWidth = signal(420);
  readonly resizingColumn = signal<'left' | 'right' | null>(null);

  private resizeStartX = 0;
  private resizeStartWidth = 0;

  async ngOnInit(): Promise<void> {
    await this.state.loadProjects();
  }

  ngAfterViewInit(): void {
    WML.Enable();
  }

  startColumnResize(side: 'left' | 'right', event: MouseEvent): void {
    event.preventDefault();
    this.resizingColumn.set(side);
    this.resizeStartX = event.clientX;
    this.resizeStartWidth = side === 'left' ? this.leftColumnWidth() : this.rightColumnWidth();

    document.addEventListener('mousemove', this.onColumnResizeMove);
    document.addEventListener('mouseup', this.onColumnResizeEnd);
  }

  private readonly onColumnResizeMove = (event: MouseEvent): void => {
    const side = this.resizingColumn();
    if (!side) return;

    const delta = event.clientX - this.resizeStartX;
    const width = Math.max(160, this.resizeStartWidth + (side === 'left' ? delta : -delta));

    if (side === 'left') {
      this.leftColumnWidth.set(width);
    } else {
      this.rightColumnWidth.set(width);
    }
  };

  private readonly onColumnResizeEnd = (): void => {
    this.resizingColumn.set(null);
    document.removeEventListener('mousemove', this.onColumnResizeMove);
    document.removeEventListener('mouseup', this.onColumnResizeEnd);
  };
}
