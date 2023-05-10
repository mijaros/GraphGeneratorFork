import {NgModule} from '@angular/core';
import {RouterModule, Routes} from '@angular/router';
import {RandomGraphComponent} from "./random-graph/random-graph.component";
import {RandomBatchComponent} from "./random-batch/random-batch.component";

const routes: Routes = [
  {path: 'graph', component: RandomGraphComponent},
  {path: 'batch', component: RandomBatchComponent},
  {path: '', redirectTo: '/graph', pathMatch: 'full'}
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule {
}
