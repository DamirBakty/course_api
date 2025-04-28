-- +goose Up

-- Create Course table
CREATE TABLE course (
                         id SERIAL PRIMARY KEY,
                         name VARCHAR(255) NOT NULL,
                         description TEXT,
                         created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                         updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

Create Chapter table
CREATE TABLE chapter (
                          id SERIAL PRIMARY KEY,
                          name VARCHAR(255) NOT NULL,
                          description TEXT,
                          "order" INT NOT NULL,
                          course_id INT NOT NULL,
                          created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                          updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                          FOREIGN KEY (course_id) REFERENCES course(id) ON DELETE CASCADE
);

-- Create Lesson table
CREATE TABLE lesson (
                         id SERIAL PRIMARY KEY,
                         name VARCHAR(255) NOT NULL,
                         description TEXT,
                         content TEXT,
                         "order" INT NOT NULL,
                         chapter_id INT NOT NULL,
                         created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                         updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                         FOREIGN KEY (chapter_id) REFERENCES chapter(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS lesson;
DROP TABLE IF EXISTS chapter;
DROP TABLE IF EXISTS course;
