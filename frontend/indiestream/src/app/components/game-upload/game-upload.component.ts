import { Component } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatProgressBar } from "@angular/material/progress-bar";
import { MatIcon } from "@angular/material/icon";
import { HttpClientModule, HttpEventType } from "@angular/common/http";
import { catchError, EMPTY, finalize, Subscription } from "rxjs";
import { NgIf } from "@angular/common";
import { GamesService } from "../../services/games.service";
import { MatFormField, MatHint, MatInput, MatLabel, MatSuffix } from "@angular/material/input";
import { FormBuilder, ReactiveFormsModule, Validators } from "@angular/forms";

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
  gameForm = this.fb.group({
    title: ['', Validators.required],
    filename: ['', Validators.required],
    file: [new DataTransfer().files, Validators.required]
  });

  constructor(private gamesService: GamesService, private fb: FormBuilder) {}

  onFileSelected(event: any) {
    const file:File = event.target.files[0];
    if (file) {
      this.gameForm.patchValue({ filename: file.name});
      this.gameForm.patchValue({ file: event.target.files});
    }
  }

  onUpload() {
    if (this.gameForm.valid) {
      const upload$ = this.gamesService.uploadGame(this.gameForm)
        .pipe(
          catchError((error) => {
            console.error('An error occured during upload', error);
            return EMPTY;
          }),
          finalize(() => this.reset())
        );

      this.uploadSub = upload$.subscribe(event => {
        if (event.type == HttpEventType.UploadProgress && event.total !== undefined) {
          this.uploadProgress = Math.round(100 * (event.loaded / event.total));
        }
      })
    }
  }

  reset() {
    this.uploadProgress = 0;
    this.uploadSub.unsubscribe();
    this.uploadSub = new Subscription();
    this.gameForm.reset();
  }
}
