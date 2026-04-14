package time

import "encoding/json"
import "time"

type UnixTime time.Time

func (ut *UnixTime) UnmarshalJSON(data []byte) error {

    var timestamp int64

    if err := json.Unmarshal(data, &timestamp); err != nil {
        return err
    }

    *ut = UnixTime(time.Unix(timestamp, 0))

    return nil

}

func (ut UnixTime) Time() time.Time {
    return time.Time(ut)
}
