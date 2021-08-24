package repo

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ozoncp/ocp-meeting-api/internal/models"
	"time"
)

var _ = Describe("Repo", func() {
	tableName := "meeting"

	now := time.Now()

	var (
		db     *sql.DB
		sqlxDB *sqlx.DB
		mock   sqlmock.Sqlmock

		ctx      context.Context
		r        Repo
		meetings []models.Meeting
	)

	BeforeEach(func() {
		var err error
		db, mock, err = sqlmock.New()
		Expect(err).Should(BeNil())
		sqlxDB = sqlx.NewDb(db, "sqlmock")

		ctx = context.Background()
		r = NewRepo(sqlxDB)

		meetings = []models.Meeting{
			{1, 1, "", now, now},
			{2, 2, "", now, now},
			{3, 3, "", now, now},
		}
	})

	AfterEach(func() {
		var err error
		mock.ExpectClose()
		err = db.Close()
		Expect(err).Should(BeNil())
	})

	Context("Test AddMany", func() {
		BeforeEach(func() {
			rows := sqlmock.NewRows([]string{"id"}).
				AddRow(1).
				AddRow(2).
				AddRow(3)
			mock.ExpectQuery("INSERT INTO "+tableName).
				WithArgs(
					meetings[0].UserId, meetings[0].Link, meetings[0].Start, meetings[0].End,
					meetings[1].UserId, meetings[1].Link, meetings[1].Start, meetings[1].End,
					meetings[2].UserId, meetings[2].Link, meetings[2].Start, meetings[2].End,
				).WillReturnRows(rows)
		})

		It("Add many meetings", func() {
			err := r.AddMany(ctx, meetings)
			Expect(err).Should(BeNil())
		})
	})

	Context("Test Add", func() {
		BeforeEach(func() {
			rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
			mock.ExpectQuery("INSERT INTO "+tableName).
				WithArgs(meetings[0].UserId, meetings[0].Link, meetings[0].Start, meetings[0].End).
				WillReturnRows(rows)

		})

		It("Add meeting", func() {
			meeting := &models.Meeting{Id: 1, UserId: 1, Link: "", Start: now, End: now}
			err := r.Add(ctx, meeting)
			Expect(err).Should(BeNil())
			Expect(meeting.Id).Should(BeEquivalentTo(1))
		})
	})

	Context("Test Describe", func() {
		BeforeEach(func() {
			rows := sqlmock.NewRows([]string{"id", "user_id", "link", "start", "\"end\""}).AddRow(
				meetings[0].Id,
				meetings[0].UserId,
				meetings[0].Link,
				meetings[0].Start,
				meetings[0].End)
			mock.ExpectQuery(
				"SELECT id, user_id, link, start, \"end\" FROM " + tableName + " WHERE").
				WithArgs(meetings[0].Id).
				WillReturnRows(rows)
		})

		It("Describe meeting", func() {
			meeting, err := r.Describe(ctx, meetings[0].Id)
			Expect(err).Should(BeNil())
			Expect(*meeting).Should(BeEquivalentTo(meetings[0]))
		})
	})

	Context("Test Update", func() {
		BeforeEach(func() {
			mock.ExpectExec("UPDATE "+tableName+" SET").
				WithArgs(
					meetings[0].UserId,
					meetings[0].Link,
					meetings[0].Start,
					meetings[0].End,
					meetings[0].Id).
				WillReturnResult(sqlmock.NewResult(1, 1))
		})

		It("Update meeting", func() {
			updated, err := r.Update(ctx, meetings[0])
			Expect(err).Should(BeNil())
			Expect(updated).Should(BeEquivalentTo(true))
		})
	})

	Context("Test List", func() {
		var limit uint64 = 3
		var offset uint64 = 1

		BeforeEach(func() {
			rows := sqlmock.NewRows([]string{"id", "user_id", "link", "start", "\"end\""}).
				AddRow(meetings[0].Id, meetings[0].UserId, meetings[0].Link, meetings[0].Start, meetings[0].End).
				AddRow(meetings[1].Id, meetings[1].UserId, meetings[1].Link, meetings[1].Start, meetings[1].End).
				AddRow(meetings[2].Id, meetings[2].UserId, meetings[2].Link, meetings[2].Start, meetings[2].End)
			mock.ExpectQuery("SELECT id, user_id, link, start, \"end\" FROM " + tableName + " LIMIT 3 OFFSET 1").
				WillReturnRows(rows)
		})

		It("List of meetings", func() {
			meetingList, err := r.List(ctx, limit, offset)
			Expect(err).Should(BeNil())
			Expect(meetingList[1].Id).Should(BeEquivalentTo(meetings[1].Id))
			Expect(meetingList[1].UserId).Should(BeEquivalentTo(meetings[1].UserId))
			Expect(meetingList[1].Link).Should(BeEquivalentTo(meetings[1].Link))
			Expect(meetingList[1].Start).Should(BeEquivalentTo(meetings[1].Start))
			Expect(meetingList[1].End).Should(BeEquivalentTo(meetings[1].End))
		})
	})

	Context("Test Remove", func() {
		BeforeEach(func() {
			query := mock.ExpectExec("DELETE FROM " + tableName + " WHERE")
			query.WithArgs(meetings[2].Id)
			query.WillReturnResult(sqlmock.NewResult(1, 1))
		})

		It("Remove meeting", func() {
			deleted, err := r.Remove(ctx, meetings[2].Id)
			Expect(err).Should(BeNil())
			Expect(deleted).Should(BeEquivalentTo(true))
		})
	})
})
