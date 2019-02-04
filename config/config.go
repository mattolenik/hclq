package config

// UseRawOutput indicates output should be raw format, not JSON
var UseRawOutput bool

// InputFile is the name of the file to read from, or null if stdin is used.
var InputFile string

// OutputFile is the name of the file to write to, or null if stdout is used.
var OutputFile string

// ModifyInPlace is set if InputFile should be edited in-place instead of outputted.
var ModifyInPlace bool

// ReplaceNTimes the max number of times a replacement can happen with set
var ReplaceNTimes int
