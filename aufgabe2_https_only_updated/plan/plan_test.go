package plan

import "testing"

const TEST_PLAN_XML = `<?xml version="1.0" encoding="UTF-8"?>
<splan>
  <tage>
    <tag><id>1</id><bezeichnung><kurz>Mo</kurz><lang>Montag</lang></bezeichnung></tag>
  </tage>
  <zeiten>
    <zeit><id>1</id><bezeichnung>8:00 - 9:15</bezeichnung></zeit>
  </zeiten>
  <mitarbeiter>
    <typ>
      <id>1</id><bezeichnung>Dozenten</bezeichnung>
      <person><id>1</id><bezeichnung><kuerzel>ahr</kuerzel><vorname>Dirk</vorname><nachname>Ahrens</nachname></bezeichnung></person>
    </typ>
  </mitarbeiter>
  <veranstaltungen>
    <veranstaltung>
      <id>1</id><bezeichnung>Test</bezeichnung><stunden>1</stunden>
      <termine><termin><tag>1</tag><zeit>1</zeit></termin></termine>
      <mitarbeiter><dozent><id>1</id></dozent></mitarbeiter>
    </veranstaltung>
  </veranstaltungen>
</splan>`

func TestParseReadsTeacherPathAndTerms(t *testing.T) {
	parsed, err := Parse([]byte(TEST_PLAN_XML))
	if err != nil {
		t.Fatalf("Parse fehlgeschlagen: %v", err)
	}
	if len(parsed.Events) != 1 {
		t.Fatalf("erwartet 1 Veranstaltung, erhalten %d", len(parsed.Events))
	}
	if len(parsed.Events[0].EmployeeIDs) != 1 || parsed.Events[0].EmployeeIDs[0] != 1 {
		t.Fatalf("Dozentenzuordnung wurde nicht gelesen: %#v", parsed.Events[0].EmployeeIDs)
	}
	if len(parsed.Events[0].Terms) != 1 {
		t.Fatalf("Termin wurde nicht gelesen: %#v", parsed.Events[0].Terms)
	}
}

func TestParseRejectsEventsWithoutTeacherAssignments(t *testing.T) {
	xmlWithoutTeacher := `<?xml version="1.0"?><splan>
	<tage><tag><id>1</id><bezeichnung><kurz>Mo</kurz><lang>Montag</lang></bezeichnung></tag></tage>
	<zeiten><zeit><id>1</id><bezeichnung>8:00 - 9:15</bezeichnung></zeit></zeiten>
	<mitarbeiter><typ><person><id>1</id><bezeichnung><kuerzel>ahr</kuerzel></bezeichnung></person></typ></mitarbeiter>
	<veranstaltungen><veranstaltung><id>1</id><stunden>1</stunden><termine><termin><tag>1</tag><zeit>1</zeit></termin></termine></veranstaltung></veranstaltungen>
</splan>`

	if _, err := Parse([]byte(xmlWithoutTeacher)); err == nil {
		t.Fatal("XML ohne Dozentenzuordnungen wurde unerwartet akzeptiert")
	}
}
