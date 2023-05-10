import {Component, Input} from '@angular/core';
import {explanation, GraphType} from "../graph-type";

@Component({
  selector: 'app-graph-kind-desc',
  templateUrl: './graph-kind-desc.component.html',
  styleUrls: ['./graph-kind-desc.component.css']
})
export class GraphKindDescComponent {
  @Input() type!: GraphType
  protected readonly explanation = explanation;
  protected readonly GraphType = GraphType;
}
