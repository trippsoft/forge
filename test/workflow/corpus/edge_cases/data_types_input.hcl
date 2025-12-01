// Various data types in input
process {
  name = "Data Types Process"
  targets = "host1"
  
  step "test" {
    name = "Test Step"
    module = "shell"
    
    input {
      string_val = "hello world"
      number_val = 42
      float_val = 3.14
      bool_val = true
      array_val = ["a", "b", "c"]
      null_val = null
      empty_string = ""
      zero_val = 0
      false_val = false
    }
  }
}
