import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {GraphRequestService} from "../graph-request.service";
import {GraphRequest} from "../graph-request";
import {GraphFormGroup} from "../graph-form-group";
import {MessageService} from "../message.service";
import {takeUntil} from "rxjs";
import {BaseComponent} from "../base.component";
import {LimitsProviderService} from "../limits-provider.service";

@Component({
  selector: 'app-graph-request',
  templateUrl: './graph-request.component.html',
  styleUrls: ['./graph-request.component.css']
})
export class GraphRequestComponent extends BaseComponent implements OnInit {

  graphForm!: GraphFormGroup

  @Input() standalone: boolean = true
  @Output() eventEmitter: EventEmitter<GraphRequest | null> = new EventEmitter<GraphRequest | null>()

  constructor(private graphService: GraphRequestService,
              private messageService: MessageService,
              private limits: LimitsProviderService) {
    super()
  }

  onSubmit() {
    let data: GraphRequest = JSON.parse(JSON.stringify(this.graphForm.value))
    this.graphService.createGraph(data).pipe(takeUntil(this.terminate)).subscribe({
      next: (value) => {
        this.graphForm.reset()
        value.deleted = new Date(value.deleted)
        this.eventEmitter.emit(value)
      }, error: (err) => {
        console.log(err)
        this.messageService.addError(err)
      }
    })
  }

  onCancel() {
    this.graphForm.reset()
    this.eventEmitter.emit(null)
  }

  ngOnInit() {
    this.graphForm = new GraphFormGroup(this.limits)
  }

}
