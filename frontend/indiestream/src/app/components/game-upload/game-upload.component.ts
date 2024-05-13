import { Component } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatProgressBar } from "@angular/material/progress-bar";
import { MatIcon } from "@angular/material/icon";
import { HttpClientModule, HttpEventType } from "@angular/common/http";
import {finalize, Subscription} from "rxjs";
import {NgIf} from "@angular/common";
import {GamesService} from "../../services/games.service";
import {MatFormField, MatHint, MatInput, MatLabel, MatSuffix} from "@angular/material/input";
import {FormBuilder, ReactiveFormsModule, Validators} from "@angular/forms";

@Component({
  selector: 'app-game-upload',
  standalone: true,
  imports: [
    MatButtonModule,
    MatProgressBar,
    MatIcon,
    NgIf,
    MatInput,
    MatHint,
    MatFormField,
    MatLabel,
    MatSuffix,
    ReactiveFormsModule,
    HttpClientModule,
  ],
  templateUrl: './game-upload.component.html',
  styleUrl: './game-upload.component.scss'
})
export class GameUploadComponent {

  uploadProgress: number = 0;
  uploadSub: Subscription = new Subscription();
  file_store: FileList = new DataTransfer().files;
  gameForm = this.fb.group({
    title: ['', Validators.required],
    filename: ['', Validators.required]
  });

  constructor(private gamesService: GamesService, private fb: FormBuilder) {}

  onFileSelected(event: any) {
    const file:File = event.target.files[0];
    if (file) {
      this.file_store = event.target.files;
      this.gameForm.patchValue({ filename: file.name});
    }
  }

  onUpload() {
    const upload$ = this.gamesService.uploadGame(this.file_store[0])
      .pipe(
        finalize(() => this.reset())
      );

    this.uploadSub = upload$.subscribe(event => {
      if (event.type == HttpEventType.UploadProgress && event.total !== undefined) {
        this.uploadProgress = Math.round(100 * (event.loaded / event.total));
      }
    })
  }

  reset() {
    this.uploadProgress = 0;
    this.uploadSub.unsubscribe();
    this.uploadSub = new Subscription();
    this.gameForm.reset();
    this.file_store = new DataTransfer().files;
  }
}
