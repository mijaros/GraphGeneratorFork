import {NgModule} from '@angular/core';
import {BrowserModule} from '@angular/platform-browser';
import {HttpClientModule} from '@angular/common/http';

import {AppRoutingModule} from './app-routing.module';
import {AppComponent} from './app.component';
import {NgbAlertModule, NgbCollapseModule, NgbModule} from '@ng-bootstrap/ng-bootstrap';
import {RandomGraphComponent} from './random-graph/random-graph.component';
import {RandomBatchComponent} from './random-batch/random-batch.component';
import {GraphRequestComponent} from './graph-request/graph-request.component';
import {ReactiveFormsModule} from '@angular/forms';
import {BatchRequestComponent} from './batch-request/batch-request.component';
import {GraphRequestFormTemplateComponent} from './graph-request-form-template/graph-request-form-template.component';
import {MessageComponent} from './message/message.component';
import {NgFor} from "@angular/common";
import {FontAwesomeModule} from '@fortawesome/angular-fontawesome';
import {GraphDetailsComponent} from './graph-details/graph-details.component';
import {LuxonModule} from "luxon-angular";
import {GraphKindDescComponent} from './graph-kind-desc/graph-kind-desc.component';

@NgModule({
  declarations: [
    AppComponent,
    RandomGraphComponent,
    RandomBatchComponent,
    GraphRequestComponent,
    BatchRequestComponent,
    GraphRequestFormTemplateComponent,
    MessageComponent,
    GraphDetailsComponent,
    GraphKindDescComponent,
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    NgbModule,
    HttpClientModule,
    ReactiveFormsModule,
    NgbAlertModule,
    NgbCollapseModule,
    NgFor,
    FontAwesomeModule,
    LuxonModule
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule {
}
