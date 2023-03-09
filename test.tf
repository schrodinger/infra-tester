variable "check_condition" {
  type    = bool
  default = true
}

variable "complex_object" {
  type = object({
    count = number
    id    = number
    str = string
    seq = list(string)
    map = object({
      key = string
    })
  })
}

resource "time_static" "example" {
  lifecycle {
    precondition {
      condition     = var.check_condition != true || var.complex_object.count == 100
      error_message = "Intented to fail"
    }
  }
}

output "sample_output" {
  value = "it's working"
}

output "another_output" {
  value = "it's working"
}

output "yet_another_output" {
  value = "it's working"
}

output "a_fourth_output" {
  value = "strings numbers 123 apple 0    orange 13567"
}
