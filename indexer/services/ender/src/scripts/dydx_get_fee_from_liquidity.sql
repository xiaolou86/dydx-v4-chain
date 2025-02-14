/**
  Returns the fee given the liquidity side.
*/
CREATE OR REPLACE FUNCTION get_fee(fill_liquidity text, event_data jsonb) RETURNS numeric AS $$
BEGIN
    IF fill_liquidity = 'TAKER' THEN
        RETURN dydx_from_jsonlib_long(event_data->'takerFee');
    ELSE
        RETURN dydx_from_jsonlib_long(event_data->'makerFee');
    END IF;
END;
$$ LANGUAGE plpgsql IMMUTABLE PARALLEL SAFE;

