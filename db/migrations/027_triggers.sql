-- Write your migrate up statements here
CREATE TRIGGER update_city before UPDATE ON city
FOR EACH ROW EXECUTE PROCEDURE update_time();
CREATE TRIGGER update_company before UPDATE ON company
FOR EACH ROW EXECUTE PROCEDURE update_time();
CREATE TRIGGER update_creation_tag before UPDATE ON creation_tag
FOR EACH ROW EXECUTE PROCEDURE update_time();
CREATE TRIGGER update_creation_type before UPDATE ON creation_type
FOR EACH ROW EXECUTE PROCEDURE update_time();
CREATE TRIGGER update_job_search_status before UPDATE ON job_search_status
FOR EACH ROW EXECUTE PROCEDURE update_time();
CREATE TRIGGER update_applicant before UPDATE ON applicant
FOR EACH ROW EXECUTE PROCEDURE update_time();
CREATE TRIGGER update_employer before UPDATE ON employer
FOR EACH ROW EXECUTE PROCEDURE update_time();
CREATE TRIGGER update_work_type before UPDATE ON work_type
FOR EACH ROW EXECUTE PROCEDURE update_time();
CREATE TRIGGER update_cv before UPDATE ON cv
FOR EACH ROW EXECUTE PROCEDURE update_time();
CREATE TRIGGER update_cv_subscriber before UPDATE ON cv_subscriber
FOR EACH ROW EXECUTE PROCEDURE update_time();
CREATE TRIGGER update_applicant_session before UPDATE ON applicant_session
FOR EACH ROW EXECUTE PROCEDURE update_time();
CREATE TRIGGER update_employer_session before UPDATE ON employer_session
FOR EACH ROW EXECUTE PROCEDURE update_time();
CREATE TRIGGER update_portfolio before UPDATE ON portfolio
FOR EACH ROW EXECUTE PROCEDURE update_time();
CREATE TRIGGER update_vacancy before UPDATE ON vacancy
FOR EACH ROW EXECUTE PROCEDURE update_time();
CREATE TRIGGER update_vacancy_subscriber before UPDATE ON vacancy_subscriber
FOR EACH ROW EXECUTE PROCEDURE update_time();
CREATE TRIGGER update_applicant_creation before UPDATE ON applicant_creation
FOR EACH ROW EXECUTE PROCEDURE update_time();
CREATE TRIGGER update_applicant_rate_to_applicant_creation before UPDATE ON applicant_rate_to_applicant_creation
FOR EACH ROW EXECUTE PROCEDURE update_time();
CREATE TRIGGER update_cv_to_portfolio before UPDATE ON cv_to_portfolio
FOR EACH ROW EXECUTE PROCEDURE update_time();
CREATE TRIGGER update_vacancy_to_creation_tag before UPDATE ON vacancy_to_creation_tag
FOR EACH ROW EXECUTE PROCEDURE update_time();
CREATE TRIGGER update_employer_rate_to_applicant_creation before UPDATE ON employer_rate_to_applicant_creation
FOR EACH ROW EXECUTE PROCEDURE update_time();
CREATE TRIGGER update_cv_to_creation_tag before UPDATE ON cv_to_creation_tag
FOR EACH ROW EXECUTE PROCEDURE update_time();
CREATE TRIGGER update_applicant_creation_to_portfolio before UPDATE ON applicant_creation_to_portfolio
FOR EACH ROW EXECUTE PROCEDURE update_time();
---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.

DROP TRIGGER update_city;
DROP TRIGGER update_company;
DROP TRIGGER update_creation_tag;
DROP TRIGGER update_creation_type;
DROP TRIGGER update_job_search_status;
DROP TRIGGER update_applicant;
DROP TRIGGER update_employer;
DROP TRIGGER update_work_type;
DROP TRIGGER update_cv;
DROP TRIGGER update_cv_subscriber;
DROP TRIGGER update_applicant_session;
DROP TRIGGER update_employer_session;
DROP TRIGGER update_portfolio;
DROP TRIGGER update_vacancy;
DROP TRIGGER update_vacancy_subscriber;
DROP TRIGGER update_applicant_creation;
DROP TRIGGER update_applicant_rate_to_applicant_creation;
DROP TRIGGER update_cv_to_portfolio;
DROP TRIGGER update_vacancy_to_creation_tag;
DROP TRIGGER update_employer_rate_to_applicant_creation;
DROP TRIGGER update_cv_to_creation_tag;
DROP TRIGGER update_applicant_creation_to_portfolio;