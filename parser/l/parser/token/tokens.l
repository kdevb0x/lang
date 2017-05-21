
data Keyword = Proc | While | Mut | Let | Func | If | Else | Return

data MathOperator = Plus | Minus | Mult | Div | Mod
data ComparisonOperator = LessOrEqual | Less | Equal | Greater | GreaterOrEqual | NotEqual
data AssignmentOperator = Assign

data Operator = MathOperator | ComparisonOperator | AssignmentOperator

data BuiltInType = Int | Bool | String
data BuiltInFunc = Print

type String string
type Whitespace string
type Unknown string
type Char string
type Token = String | Whitespace | Unknown | Char