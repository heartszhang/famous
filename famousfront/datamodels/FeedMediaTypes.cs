namespace famousfront.datamodels
{
  internal class FeedMediaTypes
  {
    public const uint FeedMediaTypeNone     = 0;
    public const uint FeedMediaTypeUnknown  = 1;
    public const uint FeedMediaTypeUrl      = 1 << 1;
    public const uint FeedMediaTypeVideo    = 1 << 2;
    public const uint FeedMediaTypeAudio    = 1 << 3;
    public const uint FeedMediaTypeImage    = 1 << 4;
    public const uint FeedMediaTypeMedia    = FeedMediaTypeAudio | FeedMediaTypeVideo;
  }
}