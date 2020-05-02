package object

type Environment struct {
	store map[string]*Object
	outer *Environment
}

func NewEnvironment() *Environment {
	s := make(map[string]*Object)
	return &Environment{store: s, outer: nil}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

// Return object reference from Environment
func (e *Environment) Get(id string) (*Object, bool) {
	obj, ok := e.store[id]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(id)
	}
	return obj, ok
}

// Set object in the environment with given ID
// Reference flag serves as indication to either deep or shallow copy the object
func (e *Environment) Set(id string, obj Object, deepCopy bool) (Object, bool) {
	objRef, found := e.Get(id) // Check if this is reassignment operation
	if found && deepCopy {     // No change to the address, only change the value
		objStored := *objRef
		if objStored.Type() != obj.Type() {
			return nil, false
		}
		// TODO: look into a different way of doing this
		switch obj.Type() {
		case ARRAY_OBJ:
			objArr, ok := (objStored).(*Array)
			if ok {
				objArr.Elements = obj.(*Array).Elements
			}
			return obj, true
		case STRING_OBJ:
			objStr, ok := objStored.(*String)
			if ok {
				objStr.Value = obj.(*String).Value
			}
			return obj, true
		case INTEGER_OBJ:
			objInt, ok := objStored.(*Integer)
			if ok {
				objInt.Value = obj.(*Integer).Value
			}
			return obj, true
		case BOOLEAN_OBJ:
			objBool, ok := objStored.(*Boolean)
			if ok {
				objBool.Value = obj.(*Boolean).Value
			}
			return obj, true
		case HASH_OBJ:
			hashObj, ok := objStored.(*Hash)
			if ok {
				hashObj.Pairs = obj.(*Hash).Pairs
			}
			return obj, true
		default:
			break
		}
	} else if !found && deepCopy {
		newObj := obj
		e.store[id] = &newObj
		return newObj, true
	}
	// Shallow-Copy
	e.store[id] = &obj
	return obj, true

}
