import { useTheme } from "@emotion/react";
import CloseOutlined from "@mui/icons-material/CloseOutlined";
import SearchOutlined from "@mui/icons-material/SearchOutlined";
import IconButton from "@mui/material/IconButton";
import InputAdornment from "@mui/material/InputAdornment";
import TextField, { type TextFieldProps } from "@mui/material/TextField";
import Tooltip from "@mui/material/Tooltip";
import { visuallyHidden } from "@mui/utils";
import type { FC } from "react";

type SearchFieldProps = Omit<
  TextFieldProps,
  "onChange" | "id" | "label" | "value" | "error"
> & {
  id: string;
  label: string;
  value: string;
  error?: string;
  onChange: (value: string) => void;
};

export const SearchField: FC<SearchFieldProps> = (props) => {
  const theme = useTheme();
  const { value, label, id, error, onChange, ...textFieldProps } = props;
  const isEmpty = value.length === 0;

  return (
    <>
      <label htmlFor={id} css={{ ...visuallyHidden }}>
        {label}
      </label>
      <TextField
        {...textFieldProps}
        error={Boolean(error)}
        helperText={error}
        type="text"
        InputProps={{
          ...textFieldProps.InputProps,
          id,
          size: "small",
          startAdornment: (
            <InputAdornment position="start">
              <SearchOutlined
                role="presentation"
                css={{
                  fontSize: 14,
                  color: theme.palette.text.secondary,
                }}
              />
            </InputAdornment>
          ),
          endAdornment: !isEmpty && (
            <Tooltip title="Clear filter">
              <IconButton
                size="small"
                onClick={() => {
                  onChange("");
                }}
              >
                <CloseOutlined
                  css={{
                    fontSize: 14,
                    color: theme.palette.text.secondary,
                  }}
                />
                <span css={{ ...visuallyHidden }}>Clear filter</span>
              </IconButton>
            </Tooltip>
          ),
        }}
        value={value}
        onChange={(e) => {
          onChange(e.currentTarget.value);
        }}
        placeholder={textFieldProps.placeholder ?? "Search..."}
        fullWidth
      />
    </>
  );
};
