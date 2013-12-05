namespace famousfront.datamodels
{
  internal class FeedStatuses
  {
    public const ulong  FeedStatusTextEmpty  =1 << 0;                  //Feed_status_text_empty uint64 = 1 << iota
    public const ulong  FeedStatusTextLittle = 1 << 1;                 //Feed_status_text_little
    public const ulong  FeedStatusTextMany   = 1 << 2;                 //Feed_status_text_many
    public const ulong  FeedStatusImageEmpty = 1 <<3;                  //Feed_status_image_empty
    public const ulong  FeedStatusImageOne   = 1 <<4;                  //Feed_status_image_one
    public const ulong  FeedStatusImageMany  = 1 <<5;                  //Feed_status_image_many
    public const ulong  FeedStatusMediaEmpty = 1 << 6;                 //Feed_status_media_empty // image, audio , video
    public const ulong  FeedStatusMediaOne   = 1 << 7;                 //Feed_status_media_one
    public const ulong  FeedStatusMediaMany  = 1 << 8;                 //Feed_status_media_many
    public const ulong  FeedStatusLinkdensityLow      = 1 << 9;             //Feed_status_linkdensity_low
    public const ulong  FeedStatusLinkdensityHigh     = 1 << 10;           //Feed_status_linkdensity_high
    public const ulong  FeedStatusFormatFlowdocument  = 1 << 11;        //Feed_status_format_flowdocument
    public const ulong  FeedStatusFormatText    = 1 << 12;              //Feed_status_format_text
    public const ulong  FeedStatusMp4           = 1 << 13;                        //Feed_status_mp4
    public const ulong  FeedStatusFlv           = 1 << 14;                        //Feed_status_flv
    public const ulong  FeedStatusContentReady  = 1 << 15;             //Feed_status_content_ready
    public const ulong  FeedStatusContentEmpty  = 1 << 16;             //Feed_status_content_empty
    public const ulong  FeedStatusContentInline = 1 << 17;             //Feed_status_content_inline
    public const ulong  FeedStatusContentExternalReady = 1 << 18;     //Feed_status_content_external_ready    
    public const ulong  FeedStatusContentExternalEmpty = 1 << 19;     //Feed_status_content_external_empty    
    public const ulong  FeedStatusContentUnresolved = 1 << 20;        //Feed_status_content_unresolved    
    public const ulong  FeedStatusContentUnavail    = 1 << 21;        //Feed_status_content_unavail    
    public const ulong  FeedStatusContentDuplicated = 1 << 22;        //Feed_status_content_duplicated    
    public const ulong  FeedStatusContentMediainline= 1 << 23;        //Feed_status_content_mediainline    
    public const ulong  FeedStatusSummaryReady      = 1 << 24;        //Feed_status_summary_ready
    public const ulong  FeedStatusSummaryEmpty      = 1 << 25;        //Feed_status_summary_empty
    public const ulong  FeedStatusSummaryInline     = 1 << 26;        //Feed_status_summary_inline
    public const ulong  FeedStatusSummaryExternalReady = 1 << 27;     //Feed_status_summary_external_ready    
    public const ulong  FeedStatusSummaryExternalEmpty = 1 << 28;     //Feed_status_summary_external_empty    
    public const ulong  FeedStatusSummaryUnresolved     = 1 << 29;    //Feed_status_summary_unresolved    
    public const ulong  FeedStatusSummaryUnavail        = 1 << 30;    //Feed_status_summary_unavail    
    public const ulong  FeedStatusSummaryDuplicated     = 1ul << 31;  //Feed_status_summary_duplicated    
    public const ulong  FeedStatusSummaryMediainline    = 1ul << 32;  //Feed_status_summary_mediainline    
  }
}