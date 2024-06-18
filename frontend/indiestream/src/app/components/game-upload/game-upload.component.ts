import { Component } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatProgressBar } from "@angular/material/progress-bar";
import { MatIcon } from "@angular/material/icon";
import { HttpClientModule, HttpEventType } from "@angular/common/http";
import { catchError, EMPTY, finalize, tap, Subscription } from "rxjs";
import { NgIf } from "@angular/common";
import { GamesService } from "../../services/games.service";
import { MatFormField, MatHint, MatInput, MatLabel, MatSuffix } from "@angular/material/input";
import { AbstractControl, FormBuilder, ReactiveFormsModule, ValidationErrors, Validators } from "@angular/forms";
import { Router } from '@angular/router';

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
  readonly allowedExtensions: string[] = ['.gba', '.gbc', '.nes', '.n64', '.v64', '.z64'];
  readonly fileNameRegex: string = '^[a-zA-Z0-9-]+$';

  uploadProgress: number = 0;
  uploadSub: Subscription = new Subscription();
  gameForm = this.fb.group({
    title: ['', {
      validators: [Validators.required, Validators.pattern(this.fileNameRegex)]
    }],
    filename: ['', { 
      validators: [Validators.required, Validators.maxLength(20), this.fileExtensionValidator()]
    }],
    file: [new DataTransfer().files, Validators.required]
  });

  constructor(private gamesService: GamesService, private fb: FormBuilder,
    private router: Router
  ) {}

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
        if (event.type === HttpEventType.UploadProgress && event.total !== undefined) {
          this.uploadProgress = Math.round(100 * (event.loaded / event.total));
        } else if (event.type === HttpEventType.Response) {
          this.router.navigate(['dashboard']);
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

  fileExtensionValidator(): Validators {
    return (control: AbstractControl): ValidationErrors | null => {
      if (!control.value) {
        return null; // Return null if there's no value (valid case)
      }
      const fileName: string = control.value;
      const isValid = this.allowedExtensions.some(ext => fileName.toLowerCase().endsWith(ext.toLowerCase()));
      return isValid ? null : { invalidExtension: { value: fileName } };
    };
  }
}
