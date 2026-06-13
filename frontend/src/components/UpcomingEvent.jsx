import React from "react";

const UpcomingEvent = ({ event }) => {
  return (
    <div className="event card">
      <div>
        <h4>{event.month}</h4>
        <strong>{event.date}</strong>
      </div>
      <div>
        <strong>{event.title}</strong>
        <p>{event?.time} - {event.location}</p>
      </div>
    </div>
  );
};

export default UpcomingEvent;
