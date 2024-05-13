import { ComponentFixture, TestBed } from '@angular/core/testing';

import { GameUploadComponent } from './game-upload.component';
import {BrowserAnimationsModule} from "@angular/platform-browser/animations";

describe('GameUploadComponent', () => {
  let component: GameUploadComponent;
  let fixture: ComponentFixture<GameUploadComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [GameUploadComponent, BrowserAnimationsModule]
    })
    .compileComponents();

    fixture = TestBed.createComponent(GameUploadComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
