import { Route } from '@angular/router';
import {GameUploadComponent} from "./components/game-upload/game-upload.component";
import {GamesOverviewComponent} from "./components/games-overview/games-overview.component";
import {AccountComponent} from "./components/account/account.component";
import {LandingPageComponent} from "./components/landing-page/landing-page.component";
import {AuthGuard} from "./guards/auth.guard";

export const routes: Route[] = [
  {
    path: '',
    component: LandingPageComponent,
  },
  {
    path: 'dashboard',
    component: GamesOverviewComponent,
    canActivate: [AuthGuard],
  },
  {
    path: 'upload',
    component: GameUploadComponent,
    canActivate: [AuthGuard],
  },
  {
    path: 'account',
    component: AccountComponent,
    canActivate: [AuthGuard],
  },
  {
    path: '**',
    redirectTo: '',
  },
];
