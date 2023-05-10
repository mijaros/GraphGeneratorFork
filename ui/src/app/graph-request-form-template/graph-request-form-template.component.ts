import {Component, Input} from '@angular/core';
import {FormGroup} from "@angular/forms";
import {explanation, GraphType, values} from "../graph-type";

@Component({
  selector: 'app-graph-request-form-template',
  templateUrl: './graph-request-form-template.component.html',
  styleUrls: ['./graph-request-form-template.component.css']
})
export class GraphRequestFormTemplateComponent {
  @Input() parentGraphForm!: FormGroup
  @Input() canHaveSeed!: boolean
  @Input() idPrefix!: string
  protected readonly explanation = explanation;
  protected readonly values = values;
  protected readonly GraphType = GraphType;

  constructor() {

  }

  get enableSeed(): boolean {
    return this.parentGraphForm.get('seed')?.enabled || false
  }

  buildId(name: string) {
    return this.idPrefix + name
  }

  seedToggle() {
    if (this.enableSeed) {
      this.parentGraphForm.get('seed')?.disable()
    } else {
      this.parentGraphForm.get('seed')?.enable()
    }
  }
}
