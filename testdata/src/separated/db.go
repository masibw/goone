package separated

import "log"

func getJob2(person Person) Job {
	var job Job
	// This is N+1 Query
	err := cnn.QueryRow("SELECT job_id, name FROM Jobs WHERE job_id = ?", person.JobID).Scan(&job.JobID, &job.Name)
	if err != nil {
		log.Fatal(err)
	}
	return job
}
