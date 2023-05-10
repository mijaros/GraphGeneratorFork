import {Component, EventEmitter, Output} from '@angular/core';
import {
  AbstractControl,
  FormBuilder,
  FormControl,
  FormGroup,
  ValidationErrors,
  ValidatorFn,
  Validators
} from "@angular/forms";
import {GraphBatch} from "../graph-batch";
import {GraphBatchServiceService} from "../graph-batch-service.service";
import {GraphFormGroup} from "../graph-form-group";
import {MessageService} from "../message.service";
import {BaseComponent} from "../base.component";
import {takeUntil} from "rxjs";
import {LimitsProviderService} from "../limits-provider.service";

@Component({
  selector: 'app-batch-request',
  templateUrl: './batch-request.component.html',
  styleUrls: ['./batch-request.component.css']
})
export class BatchRequestComponent extends BaseComponent {
  graphForm: GraphFormGroup = new GraphFormGroup(this.limits)

  batchForm!: FormGroup
  @Output() submitted: EventEmitter<GraphBatch | null> = new EventEmitter<GraphBatch | null>()

  constructor(private builder: FormBuilder,
              private batchService: GraphBatchServiceService,
              private messageService: MessageService,
              private limits: LimitsProviderService) {
    super()
    this.batchForm = new FormGroup({
      base: this.graphForm.form,
      number: new FormControl(2, {nonNullable: true, validators: Validators.required})
    }, this.numberValidator)
  }

  checkWholeNumber(n: Number): boolean {
    if (n == null) {
      return false
    }
    return n === Math.floor(n.valueOf())
  }
  numberValidator: ValidatorFn = (control: AbstractControl): ValidationErrors | null => {
    let res = {}
    let num = control.get('number')?.value
    if (num == null) {
      return {numberInvalid: true,
        numberInvalidMessage: "Number must filled and valid."}
    }
    if (!this.checkWholeNumber(num)) {
      return {numberInvalid: true, numberInvalidMessage: "Number must be whole."}
    }
    if (num <= 1) {
      res = {numberInvalid: true,
        numberInvalidMessage: "Number must be higher than 1"}
    } else if (num > this.limits.maxBatchSize) {
      res = {numberInvalid: true,
        numberInvalidMessage: `Maximum allowed size of batch is ${this.limits.maxBatchSize}.`}
    }
    return res
  }

  onSubmit() {
    let data: GraphBatch = JSON.parse(JSON.stringify(this.batchForm.getRawValue()))
    this.batchService.createBatch(data).pipe(takeUntil(this.terminate)).subscribe({
      next: (value) => {
        value.deleted = new Date(value.deleted)
        this.submitted.emit(value)
        this.batchForm.reset()
      }, error: (err) => {
        console.log(err)
        this.messageService.addError(err)
      }
    })

  }

  onCancel() {
    this.batchForm.reset()
    this.submitted.emit(null)
  }

}
