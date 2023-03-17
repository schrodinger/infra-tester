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
      error_message = "Intended to fail"
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

output "a_boolean_output" {
  value = true
}

output "a_list_output" {
  value = ["a", "b", "c"]
}

output "a_map_output" {
  value = {
    key = "value"
  }
}

output "a_number_output" {
    value = 10
}

output "a_float_output" {
    value = 123.11
}

output "a_complex_output" {
    value = {
        natural_number = 100
        float = 123.11
        str = "hello"
        seq = ["a", "b", "c"]
        map = {
            nested_map = {
                nested_key = "nested_value"
            }
            key = "value"
        }
        boolean = true
    }
}
