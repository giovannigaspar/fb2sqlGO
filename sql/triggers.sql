/* Your triggers goes here if necessary. Create them separated by a comment as
shown below. */

CREATE OR REPLACE FUNCTION TG1_FUNC() RETURNS trigger AS $$
  BEGIN
    
  END;
$$ LANGUAGE plpgsql;
/*SPLITHERE*/
CREATE TRIGGER TG1
  BEFORE INSERT OR UPDATE ON TABLE_NAME
  FOR EACH ROW
  EXECUTE PROCEDURE TG1_FUNC();
/*SPLITHERE*/

/* etc . . . */
