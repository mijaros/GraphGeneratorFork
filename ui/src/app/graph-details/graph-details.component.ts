import {Component, Input} from '@angular/core';
import {GraphRequest} from "../graph-request";
import {explanation, GraphType} from "../graph-type";
import {explain, RequestStatus} from "../request-status";
import {faCheckCircle, faCircleXmark} from "@fortawesome/free-regular-svg-icons";

@Component({
  selector: 'app-graph-details',
  templateUrl: './graph-details.component.html',
  styleUrls: ['./graph-details.component.css']
})
export class GraphDetailsComponent {

  @Input() graphDetails!: GraphRequest
  @Input() graphTemplate!: boolean


  protected readonly GraphType = GraphType;
  protected readonly Request = Request;
  protected readonly RequestStatus = RequestStatus;
  protected readonly explanation = explanation;
  protected readonly explain = explain;
  protected readonly faCheckCircle = faCheckCircle;
  protected readonly faCircleXmark = faCircleXmark;
  protected readonly Intl = Intl;
}
