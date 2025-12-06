package service

import (
	"database/sql"
	"errors"

	"github.com/dimasrizkyfebrian/coursify/internal/dto"
	"github.com/dimasrizkyfebrian/coursify/internal/model"
	"github.com/dimasrizkyfebrian/coursify/internal/repository"
)

type CourseService struct {
	Repo *repository.CourseRepository
}

func NewCourseService(repo *repository.CourseRepository) *CourseService {
	return &CourseService{Repo: repo}
}

func (s *CourseService) CreateCourse(instructorID string, req *dto.CreateCourseRequest) (*model.Course, error) {
	course := &model.Course{
		InstructorID: instructorID,
		Title:        req.Title,
		Description:  req.Description,
	}

	if req.CoverImageURL != "" {
		course.CoverImageURL = sql.NullString{String: req.CoverImageURL, Valid: true}
	} else {
		course.CoverImageURL = sql.NullString{Valid: false}
	}

	err := s.Repo.CreateCourse(course)
	if err != nil {
		return nil, err
	}

	return course, nil
}

func (s *CourseService) GetMyCourses(instructorID string) ([]model.Course, error) {
	return s.Repo.GetCoursesByInstructorID(instructorID)
}

func (s *CourseService) GetCourseByID(courseID string) (*model.Course, error) {
	return s.Repo.GetCourseByID(courseID)
}

func (s *CourseService) UpdateCourse(courseID, instructorID string, req *dto.CreateCourseRequest) error {
	existingCourse, err := s.Repo.GetCourseByID(courseID)
	if err != nil || existingCourse == nil {
		return errors.New("course not found")
	}
	if existingCourse.InstructorID != instructorID {
		return errors.New("forbidden")
	}

	existingCourse.Title = req.Title
	existingCourse.Description = req.Description
	if req.CoverImageURL != "" {
		existingCourse.CoverImageURL = sql.NullString{String: req.CoverImageURL, Valid: true}
	}

	return s.Repo.UpdateCourse(existingCourse)
}

func (s *CourseService) DeleteCourse(courseID, instructorID string) error {
	return s.Repo.DeleteCourse(courseID, instructorID)
}

func (s *CourseService) AddMaterial(courseID, instructorID string, req *dto.AddMaterialRequest) (*model.LearningMaterial, error) {
	existingCourse, err := s.Repo.GetCourseByID(courseID)
	if err != nil || existingCourse == nil {
		return nil, errors.New("course not found")
	}
	if existingCourse.InstructorID != instructorID {
		return nil, errors.New("forbidden")
	}

	material := &model.LearningMaterial{
		CourseID:    courseID,
		Title:       req.Title,
		ContentType: req.ContentType,
		TextContent: req.TextContent,
		VideoURL:    req.VideoURL,
	}

	err = s.Repo.AddMaterialToCourse(material)
	if err != nil {
		return nil, err
	}

	return material, nil
}

func (s *CourseService) GetMaterials(courseID, instructorID string) ([]model.LearningMaterial, error) {
	existingCourse, err := s.Repo.GetCourseByID(courseID)
	if err != nil || existingCourse == nil {
		return nil, errors.New("course not found")
	}
	if existingCourse.InstructorID != instructorID {
		return nil, errors.New("forbidden")
	}

	return s.Repo.GetMaterialsByCourseID(courseID)
}

func (s *CourseService) UpdateMaterial(courseID, materialID, instructorID string, req *dto.AddMaterialRequest) error {
	existingCourse, err := s.Repo.GetCourseByID(courseID)
	if err != nil || existingCourse == nil {
		return errors.New("course not found")
	}
	if existingCourse.InstructorID != instructorID {
		return errors.New("forbidden")
	}

	material := &model.LearningMaterial{
		ID:          materialID,
		CourseID:    courseID,
		Title:       req.Title,
		ContentType: req.ContentType,
		TextContent: req.TextContent,
		VideoURL:    req.VideoURL,
	}

	return s.Repo.UpdateMaterial(material)
}

func (s *CourseService) DeleteMaterial(courseID, materialID, instructorID string) error {
	existingCourse, err := s.Repo.GetCourseByID(courseID)
	if err != nil || existingCourse == nil {
		return errors.New("course not found")
	}
	if existingCourse.InstructorID != instructorID {
		return errors.New("forbidden")
	}

	return s.Repo.DeleteMaterial(courseID, materialID)
}

func (s *CourseService) GetAllCoursesPublic() ([]model.Course, error) {
	return s.Repo.GetAllCourses()
}

func (s *CourseService) EnrollStudent(studentID, courseID string) error {
	return s.Repo.EnrollStudent(studentID, courseID)
}

func (s *CourseService) GetMyEnrolledCourses(studentID string) ([]model.Course, error) {
	return s.Repo.GetEnrolledCoursesByStudentID(studentID)
}

func (s *CourseService) GetEnrolledCourseDetails(studentID, courseID string) (*dto.CourseWithMaterialsResponse, error) {
	isEnrolled, err := s.Repo.IsStudentEnrolled(studentID, courseID)
	if err != nil {
		return nil, err
	}
	if !isEnrolled {
		return nil, errors.New("forbidden")
	}

	course, err := s.Repo.GetCourseByID(courseID)
	if err != nil || course == nil {
		return nil, errors.New("course not found")
	}

	materials, err := s.Repo.GetMaterialsByCourseID(courseID)
	if err != nil {
		return nil, err
	}

	return &dto.CourseWithMaterialsResponse{
		CourseResponse: s.MapCourseToResponse(course),
		Materials:      materials,
	}, nil
}

func (s *CourseService) UploadPdfMaterial(courseID, instructorID, title, filePath string) (*model.LearningMaterial, error) {
	existingCourse, err := s.Repo.GetCourseByID(courseID)
	if err != nil || existingCourse == nil {
		return nil, errors.New("course not found")
	}
	if existingCourse.InstructorID != instructorID {
		return nil, errors.New("forbidden")
	}

	material := &model.LearningMaterial{
		CourseID: courseID,
		Title:    title,
		FileURL:  "/" + filePath,
	}

	err = s.Repo.AddFileMaterialToCourse(material)
	if err != nil {
		return nil, err
	}

	return material, nil
}

func (s *CourseService) GetEnrolledStudents(courseID, instructorID string) ([]model.User, error) {
	existingCourse, err := s.Repo.GetCourseByID(courseID)
	if err != nil || existingCourse == nil {
		return nil, errors.New("course not found")
	}
	if existingCourse.InstructorID != instructorID {
		return nil, errors.New("forbidden")
	}

	return s.Repo.GetEnrolledStudentsByCourseID(courseID)
}

// MapCourseToResponse helper
func (s *CourseService) MapCourseToResponse(course *model.Course) dto.CourseResponse {
	return dto.CourseResponse{
		ID:             course.ID,
		InstructorID:   course.InstructorID,
		InstructorName: course.InstructorName,
		Title:          course.Title,
		Description:    course.Description,
		CoverImageURL:  &course.CoverImageURL.String,
		CreatedAt:      course.CreatedAt,
		UpdatedAt:      course.UpdatedAt,
	}
}

func (s *CourseService) MapCoursesToResponse(courses []model.Course) []dto.CourseResponse {
	var response []dto.CourseResponse
	for _, course := range courses {
		response = append(response, s.MapCourseToResponse(&course))
	}
	return response
}
