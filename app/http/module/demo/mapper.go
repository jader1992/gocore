package demo

import (
	demoService "github.com/jader1992/gocore/app/provider/demo"
)

// UserModelsToUserDTOs model => dto
func UserModelsToUserDTOs(models []UserModel) []UserDto {
    var ret []UserDto
	for _, model := range models {
		t := UserDto{
			ID:   model.UserId,
			Name: model.Name,
		}
		ret = append(ret, t)
	}
	return ret
}

// StudentsToUsersDTOs Student => dto
func StudentsToUsersDTOs(students []demoService.Student) []UserDto {
    var ret []UserDto

	for _, student := range students {
		t := UserDto{
			ID:   student.ID,
			Name: student.Name,
		}
		ret = append(ret, t)
	}
	return ret
}
