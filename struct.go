package main

// type neuron struct{
//   uid int
//   title string
//   neuros []neuro
// }

// type neuro struct {
//   uid, contentType int
//   title, content string
//   dendrites []dendrite
// }

// type dendrite struct {
//   uid, count, backlinkCount int 
//   synapse []int
// }

type neuron struct{
  Uid, ContentType int
  Title, Content string
  Tags []string
  Synapse []int
}
