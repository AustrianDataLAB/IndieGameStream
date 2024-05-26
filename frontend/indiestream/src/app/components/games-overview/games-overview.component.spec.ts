import { ComponentFixture, TestBed } from '@angular/core/testing';

import { GamesOverviewComponent } from './games-overview.component';
import {BrowserAnimationsModule} from "@angular/platform-browser/animations";

describe('GamesOverviewComponent', () => {
  let component: GamesOverviewComponent;
  let fixture: ComponentFixture<GamesOverviewComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [GamesOverviewComponent, BrowserAnimationsModule]
    })
    .compileComponents();

    fixture = TestBed.createComponent(GamesOverviewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
