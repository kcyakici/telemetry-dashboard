export type DateInputProps = {
  label: string;
  date: string;
  handleDateChange: (date: string) => void;
};

export default function DateInput({
  label,
  date,
  handleDateChange,
}: DateInputProps) {
  return (
    <div>
      <label className="block text-sm mb-1">{label}</label>
      <input
        type="datetime-local"
        value={date ? new Date(date).toISOString().slice(0, 16) : ""}
        onChange={(e) => {
          const local = e.target.value; // "2019-06-24T03:21"
          const utc = new Date(local).toISOString(); // "2019-06-24T00:21:00.000Z"
          handleDateChange(utc);
        }}
        className="border rounded p-2 bg-gray-700 text-white"
      />
    </div>
  );
}
