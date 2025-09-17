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
          const local = e.target.value; // "2025-09-17T05:30"
          const isoLocal = toLocalISOString(local); // "2025-09-17T05:30:00.000Z"
          console.log(isoLocal);
          handleDateChange(isoLocal);
        }}
        className="border rounded p-2 bg-gray-700 text-white"
      />
    </div>
  );
}

const toLocalISOString = (value: string) => {
  if (!value) return "";

  const [datePart, timePart] = value.split("T");
  return `${datePart}T${timePart}:00.000Z`;
};
